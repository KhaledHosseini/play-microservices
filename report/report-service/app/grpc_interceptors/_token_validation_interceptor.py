import grpc
from grpc import ServerInterceptor
from grpc_reflection.v1alpha import reflection
import logging
import python_jwt as jwt, jwcrypto.jwk as jwk

from config import Config

class TokenValidationInterceptor(ServerInterceptor):
    def __init__(self, cfg: Config):
        self.cfg = cfg
        def noAuthHeader(ignored_request, context):
            context.abort(grpc.StatusCode.UNAUTHENTICATED, 'Missing Authorization header')
        self._noAuthHeader = grpc.unary_unary_rpc_method_handler(noAuthHeader)
        def notAuthorized(ignored_request, context):
            context.abort(grpc.StatusCode.UNAUTHENTICATED, 'Not authorized')
        self._notAuthorized = grpc.unary_unary_rpc_method_handler(notAuthorized)

    def intercept_service(self, continuation, handler_call_details):
        # Exclude reflection service from authentication so that it can retrieve the service info
        if handler_call_details.method == "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo":  
            return continuation(handler_call_details)
            
        # Extract the token from the gRPC request headers
        metadata = dict(handler_call_details.invocation_metadata)
        token = metadata.get('authorization', '')
        if token:
            logging.info(f"Token found: {token}")
            public_key = jwk.JWK.from_pem(self.cfg.AuthPublicKey)
            header, claims = jwt.verify_jwt(token, public_key, ['RS256'], checks_optional=True,ignore_not_implemented=True)
            roleStr = claims['role'].lower()
            if roleStr == "admin":
                logging.info(" role is Admin: Authorized")
                return continuation(handler_call_details)
            return self._notAuthorized
        else:
            logging.info("No token found")
            return self._noAuthHeader