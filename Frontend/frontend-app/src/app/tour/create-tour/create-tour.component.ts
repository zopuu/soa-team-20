import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, FormArray } from '@angular/forms';
import { AuthService } from 'src/app/auth/auth.service';
import { TourService } from '../tour.service';
import { TourDto } from '../tour.dto';
import { TourDifficulty } from '../tour.model';

@Component({
  selector: 'app-create-tour',
  templateUrl: './create-tour.component.html',
  styleUrls: ['./create-tour.component.css'],
})
export class CreateTourComponent {
  tourForm: FormGroup;
  difficulties = Object.values(TourDifficulty);
  transportTypes = ['Walking', 'Bicycle', 'Bus'];
  availableTags = [
    'Nature',
    'History',
    'Adventure',
    'Food',
    'Culture',
    'Relax',
  ];

  constructor(
    private fb: FormBuilder,
    private auth: AuthService,
    private tourService: TourService
  ) {
    this.tourForm = this.fb.group({
      title: ['', Validators.required],
      description: ['', [Validators.required, Validators.minLength(10)]],
      difficulty: [this.difficulties[0], Validators.required],
      transportType: [this.transportTypes[0], Validators.required],
      tags: this.fb.array(this.availableTags.map(() => false)),
      image: [''],
    });
  }

  get tagsArray(): FormArray {
    return this.tourForm.get('tags') as FormArray;
  }

  onSubmit(): void {
    if (this.tourForm.valid) {
      // Build tags list from checkboxes
      const selectedTags = this.tagsArray.controls
        .map((c, i) => (c.value ? this.availableTags[i] : null))
        .filter((t) => t) as string[];

      // Get current user id from auth service
      this.auth.whoAmI().subscribe({
        next: (me) => {
          // map frontend difficulty strings to backend numeric enum
          const diffStr: string = this.tourForm.value.difficulty;
          let difficultyNum = 0; // default Beginner
          if (diffStr === 'Begginer') difficultyNum = 0;
          else if (diffStr === 'Intermediate') difficultyNum = 1;
          else if (diffStr === 'Advanced') difficultyNum = 2;
          else if (diffStr === 'Pro') difficultyNum = 3; // Pro or unknown -> highest

          // map frontend transport type strings to backend numeric enum
          const transportStr: string = this.tourForm.value.transportType;
          let transportTypeNum = 0; // default Walking
          if (transportStr === 'Walking') transportTypeNum = 0;
          else if (transportStr === 'Bicycle') transportTypeNum = 1;
          else if (transportStr === 'Bus') transportTypeNum = 2;

          const dto: TourDto = {
            authorId: me.id as unknown as string,
            title: this.tourForm.value.title,
            description: this.tourForm.value.description,
            difficulty: difficultyNum,
            transportType: transportTypeNum,
            tags: selectedTags,
          };

          this.tourService.create(dto).subscribe({
            next: () => {
              console.log('Tour created', dto);
              // reset form and clear tag checkboxes
              this.tourForm.reset({
                difficulty: this.difficulties[0],
                transportType: this.transportTypes[0],
                image: '',
              });
              this.tagsArray.controls.forEach((c) => c.setValue(false));
            },
            error: (err) => console.error('Failed to create tour', err),
          });
        },
        error: () => {
          // redirect to login if unauthenticated
          console.error('User not authenticated');
        },
      });
    }
  }
}
