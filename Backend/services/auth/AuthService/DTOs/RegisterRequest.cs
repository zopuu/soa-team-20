using System.ComponentModel.DataAnnotations;

namespace AuthService.DTOs {
    public class RegisterRequest {
        [Required, MinLength(3),MaxLength(20)]
        public string Username { get; set; }

        [Required, EmailAddress]
        public string Email { get; set; }

        [Required]
        public string Password { get; set; }

        [Required]
        [RegularExpression("Guide|Tourist", ErrorMessage = "Role must be Guide or Tourist.")]
        public string Role { get; set; }

    }
}
