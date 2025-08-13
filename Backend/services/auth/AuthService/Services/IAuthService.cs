using AuthService.DTOs;
using AuthService.Models;

namespace AuthService.Services {
    public interface IAuthService
    {
        Task<User> RegisterAsync(RegisterRequest dto);
        Task<String> LoginAsync(LoginRequest dto);
    }
}
