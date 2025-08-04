using AuthService.Data;
using AuthService.Models;
using Microsoft.EntityFrameworkCore;

namespace AuthService.Repositories {
    public class UserRepository : IUserRepository
    {
        private readonly AppDbContext _db;
        public UserRepository(AppDbContext db) => _db = db;
        public Task<User> GetByUsernameAsync(string username) =>
            _db.Users.SingleOrDefaultAsync(u => u.Username == username);
        
        public Task<User> GetByIdAsync(int id) =>
            _db.Users.SingleOrDefaultAsync(u => u.Id == id);

        public Task<bool> UsernameExistsAsync(string username) =>
            _db.Users.AnyAsync(u => u.Username == username);

        public async Task AddAsync(User user)
        {
            await _db.Users.AddAsync(user);
        }

        public Task SaveChangesAsync() => _db.SaveChangesAsync();

    }
}
