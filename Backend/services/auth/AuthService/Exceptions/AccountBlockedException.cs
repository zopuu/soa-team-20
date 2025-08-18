namespace AuthService.Exceptions {
    public class AccountBlockedException : Exception{
        public AccountBlockedException(string message = "Account is blocked") : base(message) { }
    }
}
