using System.Security.Claims;
using AuthService.DTOs;
using AuthService.Models;
using AuthService.Repositories;
using AuthService.Utills;
using BCrypt.Net;

namespace AuthService.Services {
    public class AuthService : IAuthService
    {
        private readonly IUserRepository _repo;
        private readonly JwtUtils _jwt;
        private readonly IConfiguration _cfg;

        public AuthService(
            IUserRepository repo,
            JwtUtils jwt,
            IConfiguration cfg)
        {
            _cfg = cfg;
            _repo = repo;
            _jwt = jwt;
        }

        public async Task<User> RegisterAsync(RegisterRequest dto)
        {
            if (await _repo.UsernameExistsAsync(dto.Username))
                throw new InvalidOperationException("Username taken.");
            var user = new User
            {
                Username = dto.Username,
                Email = dto.Email,
                Role = dto.Role,
                PasswordHash = BCrypt.Net.BCrypt.HashPassword(dto.Password)
            };
            await _repo.AddAsync(user);
            await _repo.SaveChangesAsync();
            return user;
        }

        public async Task<string> LoginAsync(LoginRequest dto)
        {
            var user = await _repo.GetByUsernameAsync(dto.Username);
            // 1) Handle “user not found”
            if (user == null)
                throw new UnauthorizedAccessException("Invalid credentials.");

            // 2) Check the password once
            var isValid = BCrypt.Net.BCrypt.Verify(dto.Password, user.PasswordHash);

            // 3) Throw only if the password is *not* valid
            if (!isValid)
                throw new UnauthorizedAccessException("Invalid credentials.");

            var claims = new List<Claim>()
            {
                new Claim(ClaimTypes.NameIdentifier, user.Id.ToString()),
                new Claim(ClaimTypes.Name, user.Username),
                new Claim(ClaimTypes.Role, user.Role),
                new Claim(ClaimTypes.Email, user.Email)
            };
            return _jwt.GenerateToken(claims);

        }
    }
}
