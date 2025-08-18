import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../auth.service';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import Swal from 'sweetalert2';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  constructor(
    private router: Router,
    private authService: AuthService
  ) {}
  loginForm = new FormGroup({
    username: new FormControl('', {
      validators: [Validators.required],
      updateOn: 'blur'
    }),
    password: new FormControl('', [Validators.required])
  });

  login() {
    if(this.loginForm.invalid) return;

    const { username, password } = this.loginForm.value as { username: string; password: string };
    this.authService.login({username, password}).subscribe({
      next: () => {
        this.authService.whoAmI().subscribe(user =>{

          Swal.fire({
            icon: 'success',
            title: 'Login successful',
            text: 'You have successfully logged in.',
            showConfirmButton: false,
            timer: 2500
          });
          if(user.role === 'Admin') {
            this.router.navigate(['/admin']);
          } else {
            this.router.navigate(['/']);
          }
        });
      },
      error: () =>{
        this.loginForm.setErrors({ invalidCredentials: true })
      } 
    });
  }
}
