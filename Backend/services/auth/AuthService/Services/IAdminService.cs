using AuthService.DTOs;
using AuthService.Models;

namespace AuthService.Services {
    public interface IAdminService {
        Task<IEnumerable<UserDto>> GetAllAsync(CancellationToken ct = default);
        Task BlockAsync(int id);
        Task UnblockAsync(int id);
    }
}
