using System.Security.Claims;
using AuthService.DTOs;
using AuthService.Exceptions;
using AuthService.Services;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using System.IdentityModel.Tokens.Jwt;

namespace AuthService.Controllers {
    [ApiController]
    [Route("api/auth")]
    public class AuthController : ControllerBase {
        private readonly IAuthService _authService;
        public AuthController(IAuthService authService) {
            _authService = authService;
        }

        [HttpPost("register")]
        public async Task<IActionResult> Register(RegisterRequest dto)
        {
            try
            {
                var user = await _authService.RegisterAsync(dto);
                return Ok(new
                {
                    user.Id,
                    user.Username,
                    user.Email,
                    user.Role
                });
            }
            catch (InvalidOperationException ex)
            {
                return BadRequest(ex.Message);
            }
        }

        [HttpPost("login")]
        public async Task<IActionResult> Login(LoginRequest dto)
        {
            try {
                var token = await _authService.LoginAsync(dto);
                return Ok(new { token });
            }
            catch (UnauthorizedAccessException ex) {
                return Unauthorized("Invalid username or password.");
            }
            catch (AccountBlockedException ex) {
                return StatusCode(StatusCodes.Status423Locked, new { message = ex.Message });
            }
        }

        [Authorize]
        [HttpGet("whoami")]
        public IActionResult WhoAmI()
        {
            var id =
                User.FindFirstValue(ClaimTypes.NameIdentifier) ??
                User.FindFirstValue(JwtRegisteredClaimNames.Sub) ??
                User.FindFirst("uid")?.Value ??
                Request.Headers["X-User-Id"].FirstOrDefault(); // gateway fallback

            var username = User.Identity?.Name ?? User.FindFirstValue(ClaimTypes.Name);
            var role     = User.FindFirstValue(ClaimTypes.Role) ?? User.FindFirst("role")?.Value;
            var email    = User.FindFirstValue(ClaimTypes.Email) ?? User.FindFirst("email")?.Value;

            return Ok(new { id, username, role, email });
        }

    }
}
