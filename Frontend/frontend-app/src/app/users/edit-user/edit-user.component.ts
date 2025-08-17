// src/app/users/edit-user/edit-user.component.ts
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { UserService, User, UpdateUserPayload } from '../../services/user.service';

@Component({
  selector: 'app-edit-user',
  templateUrl: './edit-user.component.html',
  styleUrls: ['./edit-user.component.css']
})
export class EditUserComponent implements OnInit {
  id!: number;
  form!: FormGroup;
  loading = true;
  saving = false;
  error = '';

  constructor(
    private route: ActivatedRoute,
    private fb: FormBuilder,
    private users: UserService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.id = Number(this.route.snapshot.paramMap.get('id'));

    this.form = this.fb.group({
      FirstName: ['', [Validators.required, Validators.maxLength(50)]],
      LastName:  ['', [Validators.required, Validators.maxLength(50)]],
      ProfilePhoto: [''],
      Description: [''],
      Moto: ['']
    });

    // Prefill with current values
    this.users.getById(this.id).subscribe({
      next: (u: User) => {
        this.form.patchValue({
          FirstName: u.FirstName ?? '',
          LastName: u.LastName ?? '',
          ProfilePhoto: u.ProfilePhoto ?? '',
          Description: u.Description ?? '',
          Moto: u.Moto ?? ''
        });
        this.loading = false;
      },
      error: () => {
        this.error = 'Failed to load user.';
        this.loading = false;
      }
    });
  }

  save(): void {
    if (this.form.invalid) return;
    this.saving = true;

    const payload: UpdateUserPayload = this.form.value;

    this.users.updateUser(this.id, payload).subscribe({
      next: () => {
        this.saving = false;
        this.router.navigate(['/users', this.id, 'view']);
      },
      error: () => {
        this.error = 'Failed to save changes.';
        this.saving = false;
      }
    });
  }

  cancel(): void {
    this.router.navigate(['/users', this.id, 'view']);
  }
}
