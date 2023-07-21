use crate::proto::{RoleType, CreateUserRequest,LoginUserRequest, RefreshTokenRequest, LogOutRequest};
    use crate::models::user::db::Validate;
    #[test]
    fn create_user_request_validation() {
        let mut request = CreateUserRequest {
            name:"".into(),
            email:"".into(),
            password:"".into(),
            role: RoleType::RoletypeAdmin.into()
        };
        let mut result = request.validate();
        assert!(result.is_err());
        
        request.name = "name".into();
        result = request.validate();
        assert!(result.is_err());

        request.email = "example@example.com".into();
        result = request.validate();
        assert!(result.is_err());

        request.password = "password".into();
        result = request.validate();
        assert!(!result.is_err());

    }
    #[test]
    fn login_user_request_validation() {
        let mut request = LoginUserRequest {
            email:"".into(),
            password:"".into()
        };
        let mut result = request.validate();
        assert!(result.is_err());

        request.email = "example@example.com".into();
        result = request.validate();
        assert!(result.is_err());

        request.password = "password".into();
        result = request.validate();
        assert!(!result.is_err());

    }
    #[test]
    fn refresh_token_request_validation() {
        let mut request = RefreshTokenRequest {
            refresh_token:"".into(),
        };
        let mut result = request.validate();
        assert!(result.is_err());

        request.refresh_token = "token".into();
        result = request.validate();
        assert!(!result.is_err());

    }

    #[test]
    fn logout_user_request_validation() {
        let mut request = LogOutRequest {
            refresh_token:"".into(),
            access_token:"".into()
        };
        let mut result = request.validate();
        assert!(result.is_err());

        request.refresh_token = "token".into();
        result = request.validate();
        assert!(result.is_err());

        request.access_token = "token".into();
        result = request.validate();
        assert!(!result.is_err());

    }