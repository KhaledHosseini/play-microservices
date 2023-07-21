use chrono::{DateTime, NaiveDateTime, Utc};

use hyper::{ Request, Response};
use hyper::header::{HeaderName, HeaderValue};
use std::task::{Context, Poll};
use tonic::{
    body::BoxBody,
    transport::{Body, NamedService},
    Status,
};
use tower::Service;

use crate::utils::jwtoken;
use crate::config::Config;
use log::info;

#[derive(Debug, Clone)]
pub(crate) struct AuthenticatedService<S> {
    pub(crate) inner: S,
    pub env: Config,
}

impl<S: NamedService> NamedService for AuthenticatedService<S> {
    const NAME: &'static str = S::NAME;
}

impl<S> Service<Request<Body>> for AuthenticatedService<S>
where
    S: Service<Request<Body>, Response = Response<BoxBody>> + NamedService + Clone + Send + 'static,
    S::Future: Send + 'static,
{
    type Response = S::Response;
    type Error = S::Error;
    type Future = futures::future::BoxFuture<'static, Result<Self::Response, Self::Error>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx)
    }

    fn call(&mut self, req: Request<Body>) -> Self::Future {
        // This is necessary because tonic internally uses `tower::buffer::Buffer`.
        // See https://github.com/tower-rs/tower/issues/547#issuecomment-767629149
        // for details on why this is necessary
        let clone = self.inner.clone();
        let mut inner = std::mem::replace(&mut self.inner, clone);
        let public_key = self.env.access_token_public_key.to_owned();
        Box::pin(async move {
            match authenticate(req, public_key) {
                Ok(req) => inner.call(req).await,
                Err(e) => match e.code() {
                    tonic::Code::Unauthenticated | tonic::Code::PermissionDenied => Ok(e.to_http()),
                    _ => Ok(Status::internal("Unable to process with authorization").to_http()),
                },
            }
        })
    }
}

fn authenticate(
    request: Request<Body>,
    public_key: Vec<u8>,
) -> Result<Request<Body>, Status> {
    info!("authenticate:start {:#?}", &request);
    let headers = request.headers();
    let access_token = match headers.get("authorization") {
        Some(token) => match token.to_str() {
            Ok(token) => token,
            Err(_) => return Err(Status::unauthenticated("Token malformed (string)")),
        },
        None => return Err(Status::unauthenticated("Token not found")),
    };
    info!("authenticate:  found access token{:#?}", &access_token);
    let access_token_details =
    match jwtoken::verify_jwt_token(public_key, access_token)
    {
        Ok(token_details) => token_details,
        Err(e) => {
            return Err(Status::unauthenticated(format!("Token validation failed with error: {}", e)));
        }
    };
    info!("authenticate:  access token verified. {:#?}", &access_token_details);
    //check expiration 
    match access_token_details.expires_in {
        Some(exp)=> {
            if let Some(exp) = NaiveDateTime::from_timestamp_opt(exp, 0) {
                let expires = DateTime::<Utc>::from_utc(exp, Utc);
                if Utc::now() > expires {
                    return Err(Status::unauthenticated("Token has expired"));
                }
            }else {
                return Err(Status::unauthenticated("Error checking token expiry"));
            }
        },
        None => {
            return Err(Status::unauthenticated("Error checking token expiry"));
        }
    }
    info!("authenticate:  set auth headers. ");
    //add authorization values to request
    let mut modified_request = request;
    modified_request.headers_mut().insert(
        HeaderName::from_static("user_role"),
        HeaderValue::from_str(&access_token_details.role).map_err(|_| Status::internal("Error setting role"))?,
    );
    modified_request.headers_mut().insert(
        HeaderName::from_static("user_id"),
        HeaderValue::from_str(&access_token_details.user_id.to_string()).map_err(|_| Status::internal("Error setting user id"))?,
    );

    Ok(modified_request)
}
