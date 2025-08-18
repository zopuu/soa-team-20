using AuthService.Models;

namespace AuthService.Repositories {
    public interface IUserRepository
    {
        Task<User> GetByUsernameAsync(string username);
        Task<User> GetByIdAsync(int id);
        Task AddAsync(User user);
        Task<bool> UsernameExistsAsync(string username);
        Task SaveChangesAsync();
        Task<List<User>> GetAllAsync(CancellationToken ct = default);
    }
}
