using AuthService.Models;
using Microsoft.EntityFrameworkCore;
namespace AuthService.Data {
    public class AppDbContext : DbContext{
        public AppDbContext(DbContextOptions<AppDbContext> opts) : base(opts) { }
        public DbSet<User> Users { get; set; }
    }
}
