export class AuthError extends Error {
    constructor(message='Not Authorized') {
        super(message)
        this.name = 'AuthError'
    }
}