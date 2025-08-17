using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace AuthService.Models 
{
    public class User
    {
        [Key]
        [DatabaseGenerated(DatabaseGeneratedOption.Identity)]
        public int Id { get; set; }

        [Required, MaxLength(20)] public string Username { get; set; }

        [Required, EmailAddress] public string Email { get; set; }

        [Required] public string PasswordHash { get; set; }

        [Required]
        [RegularExpression("Guide|Tourist|Admin", ErrorMessage = "Role must be Guide, Tourist or Admin")]
        public string Role { get; set; }
        public string FirstName { get; set; }
        public string LastName { get; set; }
        public string Description { get; set; }
        public string Moto { get; set; }
        public string ProfilePhoto { get; set; }
        public DateTimeOffset CreatedAt { get; set; } = DateTimeOffset.UtcNow;
    }
}
