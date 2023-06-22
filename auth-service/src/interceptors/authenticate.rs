use crate::token;

use chrono::{DateTime, Duration, NaiveDateTime, Utc};
use hyper::Body;
use std::{
    task::{Context, Poll},
    time::Duration,
};
use tonic::{body::BoxBody, transport::Server, Request, Response, Status};
use tower::{Layer, Service};

use redis::{Client, AsyncCommands};
use crate::token;
use crate::config::Config;

#[derive(Debug, Clone)]
pub(crate) struct AuthenticateMiddleware<S> {
    pub(crate) inner: S,
    redis_client: Client,
    env: Config,
}

impl<S> Service<hyper::Request<Body>> for AuthenticateMiddleware<S>
where
    S: Service<hyper::Request<Body>, Response = hyper::Response<BoxBody>> + Clone + Send + 'static,
    S::Future: Send + 'static,
{
    type Response = S::Response;
    type Error = S::Error;
    type Future = futures::future::BoxFuture<'static, Result<Self::Response, Self::Error>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx)
    }

    fn call(&mut self, req: hyper::Request<Body>) -> Self::Future {
        // This is necessary because tonic internally uses `tower::buffer::Buffer`.
        // See https://github.com/tower-rs/tower/issues/547#issuecomment-767629149
        // for details on why this is necessary
        let clone = self.inner.clone();
        let mut inner = std::mem::replace(&mut self.inner, clone);
        // Get the headers from the request

        Box::pin(async move {
            let headers = req.headers();
            // Retrieve the authorization header value from the headers
            let access_token = match headers.get("Authorization") {
                Some(value) => value.to_str().map_err(|_| {
                    // Handle invalid header value
                    let access_token_details =
                    match token::verify_jwt_token(self.context.env.refresh_token_public_key.to_owned(), &value)
                    {
                        Ok(token_details) => token_details,
                        Err(e) => {
                            return Err(Status::internal(format!("Token validation failed with error: {}", e)));
                        }
                    };
                    
                    //check expiration 
                    let expires = DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(access_token_details.expires_in, 0), Utc);
                    if Utc::now() > expires {
                        return Err(Status::unauthenticated("Token has expired"));
                    }
                })?,
                None => {
                    // Handle missing header
                    return Err(Status::internal("No Auth header found"));
                }
            };

            let response = inner.call(req).await?;

            Ok(response)
        })
    }
}