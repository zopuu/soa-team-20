using AuthService.Models;

namespace AuthService.DTOs {
    public class UserDto {
        public int Id { get; set; }
        public string Username { get; set; }
        public string Email { get; set; }
        public string Role { get; set; }
        public string? FirstName { get; set; }
        public string? LastName { get; set; }
        public string? Description { get; set; }
        public string? Moto { get; set; }
        public string? ProfilePhoto { get; set; }
        public UserStatus Status { get; set; }
        public DateTimeOffset CreatedAt { get; set; }
        public DateTimeOffset? BlockedAt { get; set; }
    }

}
