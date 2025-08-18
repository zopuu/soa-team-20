using AuthService.DTOs;
using AuthService.Services;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;

namespace AuthService.Controllers {
    [ApiController]
    [Route("api/admin")]
    [Authorize(Roles = "Admin")]
    public class AdminController : ControllerBase {
        private readonly IAdminService _svc;
        public AdminController(IAdminService svc) => _svc = svc;
        [HttpGet]
        public async Task<ActionResult<IEnumerable<UserDto>>> GetAll(CancellationToken ct) {
            var users = await _svc.GetAllAsync(ct);
            return Ok(users);
        }
        [HttpPut("{id:int}/block")]
        public async Task<IActionResult> Block(int id) {
            try {
                await _svc.BlockAsync(id);
                return NoContent();
            }
            catch(KeyNotFoundException) { return NotFound(); }
            catch(InvalidOperationException ex) { return BadRequest(ex.Message); }
        }
        [HttpPut("{id:int}/unblock")]
        public async Task<IActionResult> Unblock(int id) {
            try {
                await _svc.UnblockAsync(id);
                return NoContent();
            }
            catch (KeyNotFoundException) { return NotFound(); }
            catch (InvalidOperationException ex) { return BadRequest(ex.Message); }
        }

    }
}
