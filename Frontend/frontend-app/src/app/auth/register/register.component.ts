import { Component } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, ValidationErrors, Validators } from '@angular/forms';
import { AuthService } from '../auth.service';
import { Router } from '@angular/router';
import Swal from 'sweetalert2';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent {
  registrationForm = new FormGroup({
    username:       new FormControl('', Validators.required),
    email:          new FormControl('', [Validators.required, Validators.email]),
    role:           new FormControl('Turista', Validators.required),
    password:       new FormControl('', [Validators.required, Validators.minLength(8)]),
    confirmPassword:new FormControl('', Validators.required)
  }, { validators: this.passwordMatchValidator });
  constructor(
    private authService: AuthService,
    private router: Router
  ) {}
  private passwordMatchValidator(group: AbstractControl): ValidationErrors | null {
    const pass  = group.get('password')?.value;
    const conf  = group.get('confirmPassword')?.value;
  
    if (pass && conf && pass !== conf) {
      // mark the child
      group.get('confirmPassword')?.setErrors({ passwordMismatch: true });
      return { passwordMismatch: true };
    } else {
      // clear mismatch error but keep any others (e.g. required/minLength)
      const errors = group.get('confirmPassword')?.errors;
      if (errors) {
        delete errors['passwordMismatch'];
        if (!Object.keys(errors).length) {
          group.get('confirmPassword')?.setErrors(null);
        }
      }
      return null;
    }
  }
  register() {
    if (this.registrationForm.invalid) return;
    const { username, email, role, password } =
      this.registrationForm.value as {
        username: string;
        email:    string;
        role:     string;
        password: string;
        confirmPassword: string;
      };
    this.authService.register({ username, email, role, password }).subscribe({
        next: () => {
          Swal.fire({
            icon: 'success',
            title: 'Registration successful',
            showConfirmButton: false,
            timer: 2000
          });
          this.router.navigate(['/auth/login']);
        },
        error: (err) => {
          // assume backend returns 400 with message "Username taken"
          if (err.status === 400 && err.error?.includes('Username')) {
            this.registrationForm.get('username')?.setErrors({ usernameTaken: true });
          } else {
            Swal.fire('Error', 'An unexpected error occurred.', 'error');
          }
        }
      });
  }
}
