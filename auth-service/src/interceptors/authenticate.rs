use chrono::{DateTime, NaiveDateTime, Utc};

use hyper::{ Request, Response};
use std::task::{Context, Poll};
use tonic::{
    body::BoxBody,
    transport::{Body, NamedService},
    Status,
};
use tower::{ Service};

use crate::token;
use crate::config::Config;

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
    public_key: String,
) -> Result<Request<Body>, Status> {
    
    let access_token = match request.headers().get("authorization") {
        Some(token) => match token.to_str() {
            Ok(token) => token,
            Err(_) => return Err(Status::unauthenticated("Token malformed (string)")),
        },
        None => return Err(Status::unauthenticated("Token not found")),
    };
    let access_token_details =
    match token::verify_jwt_token(public_key, access_token)
    {
        Ok(token_details) => token_details,
        Err(e) => {
            return Err(Status::internal(format!("Token validation failed with error: {}", e)));
        }
    };
    //check expiration 
    match access_token_details.expires_in {
        Some(exp)=> {
            if let Some(exp) = NaiveDateTime::from_timestamp_opt(exp, 0) {
                let expires = DateTime::<Utc>::from_utc(exp, Utc);
                if Utc::now() > expires {
                    return Err(Status::unauthenticated("Token has expired"));
                }
            }else {
                return Err(Status::internal("Error checking token expiry"));
            }
        },
        None => {
            return Err(Status::internal("Error checking token expiry"));
        }
    }
    
    // if claims.issuer_uid.len() == 0 {
    //     return Err(Status::unauthenticated("Issuer UID empty"));
    // }
    // let uid = match HeaderValue::from_str(&claims.issuer_uid) {
    //     Ok(uid) => uid,
    //     Err(_) => return Err(Status::unknown("Internal service error")),
    // };
    // request.headers_mut().insert(METADATA_USER_ID, uid);
    Ok(request)
}
