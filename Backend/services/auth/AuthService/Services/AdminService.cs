using AuthService.DTOs;
using AuthService.Repositories;

namespace AuthService.Services {
    public class AdminService : IAdminService {
        private readonly IUserRepository _repo;
        public AdminService(IUserRepository repo) => _repo = repo;

        public async Task<IEnumerable<UserDto>> GetAllAsync(CancellationToken ct = default) {
            var items = await _repo.GetAllAsync(ct);
            return items.Select(u => new UserDto {
                Id = u.Id,
                Username = u.Username,
                Email = u.Email,
                Role = u.Role,
                FirstName = u.FirstName,
                LastName = u.LastName,
                Description = u.Description,
                Moto = u.Moto,
                ProfilePhoto = u.ProfilePhoto,
                Status = u.Status,
                CreatedAt = u.CreatedAt,
                BlockedAt = u.BlockedAt
            });
        }
        public async Task BlockAsync(int id) {
            var user = await _repo.GetByIdAsync(id) ?? throw new KeyNotFoundException("User not found.");
            if (user.Status == Models.UserStatus.BLOCKED)
                throw new InvalidOperationException("User already blocked.");
            user.Status = Models.UserStatus.BLOCKED;
            user.BlockedAt = DateTimeOffset.UtcNow;
            await _repo.SaveChangesAsync();            
        }
        public async Task UnblockAsync(int id) {
            var user = await _repo.GetByIdAsync(id) ?? throw new KeyNotFoundException("User not found.");
            if (user.Status == Models.UserStatus.ACTIVE)
                throw new InvalidOperationException("User already active.");
            user.Status = Models.UserStatus.ACTIVE;
            user.BlockedAt = null;
            await _repo.SaveChangesAsync();
        }

    }
}
