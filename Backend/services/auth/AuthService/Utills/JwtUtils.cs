using Microsoft.IdentityModel.Tokens;
using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Text;

namespace AuthService.Utills {
    public class JwtUtils
    {
        private readonly string _key;
        private readonly string _issuer;
        private readonly string _audience;

        public JwtUtils(IConfiguration cfg)
        {
            _key = cfg["Jwt:Key"];
            _issuer = cfg["Jwt:Issuer"];
            _audience = cfg["Jwt:Audience"];
        }

        public string GenerateToken(IEnumerable<Claim> claims)
        {
            var key = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(_key));
            var creds = new SigningCredentials(key, SecurityAlgorithms.HmacSha256);
            var token = new JwtSecurityToken(
                issuer: _issuer,
                audience: _audience,
                claims: claims,
                expires: DateTime.UtcNow.AddDays(1),
                signingCredentials: creds
            );
            return new JwtSecurityTokenHandler().WriteToken(token);
        }
    }
}
