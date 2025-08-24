import {
  Component,
  AfterViewInit,
  OnDestroy,
  ViewChild,
  ElementRef,
  OnInit,
} from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import * as L from 'leaflet';
import { KeypointService } from '../keypoint.service';
import { KeyPointDto } from '../keypoint.dto';
import { AuthService } from 'src/app/auth/auth.service';
@Component({
  selector: 'app-create-keypoint',
  templateUrl: './create-keypoint.component.html',
  styleUrls: ['./create-keypoint.component.css'],
})
export class CreateKeypointComponent
  implements OnInit, AfterViewInit, OnDestroy
{
  @ViewChild('mapContainer', { static: true })
  mapContainer!: ElementRef<HTMLDivElement>;
  private map?: L.Map;
  private marker?: L.Marker; // marker for the point being created/edited
  private existingMarkers: L.Marker[] = []; // markers for already-saved keypoints
  private existingIcon: L.Icon = L.icon({
    iconUrl: 'assets/icons/keypointmark.svg',
    iconSize: [32, 32],
    iconAnchor: [16, 32],
  });
  private positionIcon: L.Icon = L.icon({
    iconUrl: 'assets/icons/positionmark.svg',
    iconSize: [32, 32],
    iconAnchor: [16, 32],
  });

  form: FormGroup;
  tourId?: string;
  userId?: string;
  constructor(
    private fb: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private keypointService: KeypointService,
    private auth: AuthService
  ) {
    this.form = this.fb.group({
      tourId: ['', Validators.required],
      name: ['', Validators.required],
      description: [''],
      image: [''],
      latitude: [null, Validators.required],
      longitude: [null, Validators.required],
    });
    this.auth.whoAmI().subscribe({
      next: (me) => {
        this.userId = me.id?.toString();
      },
      error: () => {
        this.userId = undefined;
      },
    });
  }

  cancel(): void {
    if (this.tourId) {
      this.router.navigate(['users', this.userId, 'tours']);
    } else {
      this.router.navigate(['/tours']);
    }
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe((pm) => {
      const id = pm.get('id');
      if (id) {
        this.tourId = id;
        this.form.patchValue({ tourId: id });
        // if map already initialized, load keypoints now
        if (this.map) {
          this.loadKeypointsForTour();
        }
      }
    });
  }

  ngAfterViewInit(): void {
    // initialize a basic Leaflet map that supports zoom in/out
    this.map = L.map(this.mapContainer.nativeElement, {
      center: [20, 0],
      zoom: 2,
      zoomControl: true,
    });

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors',
      maxZoom: 19,
    }).addTo(this.map);

    // listen for clicks to set coordinates
    this.map.on('click', (e: L.LeafletMouseEvent) => {
      const lat = e.latlng.lat;
      const lng = e.latlng.lng;
      this.setMarkerAndCoords(lat, lng);
    });

    // small delay to ensure correct sizing in complex layouts
    setTimeout(() => this.map?.invalidateSize(), 200);

    // load existing keypoints for this tour (if tourId present)
    this.loadKeypointsForTour();
  }

  private setMarkerAndCoords(lat: number, lng: number) {
    // place or move marker
    if (!this.map) return;
    if (this.marker) {
      this.marker.setLatLng([lat, lng]);
    } else {
      this.marker = L.marker([lat, lng], { icon: this.positionIcon }).addTo(
        this.map
      );
    }
    // update form fields
    this.form.patchValue({ latitude: lat, longitude: lng });
  }

  private loadKeypointsForTour(): void {
    // clear previous markers
    this.existingMarkers.forEach((m) => m.remove());
    this.existingMarkers = [];
    if (!this.tourId || !this.map) return;

    this.keypointService.getByTour(this.tourId).subscribe({
      next: (kps) => {
        (kps || []).forEach((kp: any) => {
          try {
            const lat = kp.coordinates?.latitude;
            const lng = kp.coordinates?.longitude;
            if (lat == null || lng == null) return;
            const m = L.marker([lat, lng], { icon: this.existingIcon }).addTo(
              this.map!
            );
            const title = kp.name || 'Keypoint';
            const desc = kp.description ? `<div>${kp.description}</div>` : '';
            m.bindPopup(`<b>${title}</b>${desc}`);
            this.existingMarkers.push(m);
          } catch (e) {
            console.warn('Failed to render keypoint', kp, e);
          }
        });
      },
      error: (err) => console.error('Failed to load keypoints for tour', err),
    });
  }

  ngOnDestroy(): void {
    if (this.map) {
      this.map.remove();
      this.map = undefined;
    }
  }

  submit(): void {
    if (this.form.invalid) return;
    const v = this.form.value;
    const dto: KeyPointDto = {
      tourId: v.tourId,
      name: v.name,
      description: v.description,
      image: v.image,
      coordinates: {
        latitude: Number(v.latitude),
        longitude: Number(v.longitude),
      },
    };

    this.keypointService.create(dto).subscribe({
      next: (kp) => {
        // Clear input fields but keep tourId so user can add another keypoint
        this.form.patchValue({
          name: '',
          description: '',
          image: '',
          latitude: null,
          longitude: null,
        });
        this.form.markAsPristine();
        this.form.markAsUntouched();

        // remove the temporary create marker from the map
        if (this.marker) {
          this.marker.remove();
          this.marker = undefined;
        }

        // reload existing markers so the newly created keypoint appears
        this.loadKeypointsForTour();
      },
      error: (err) => {
        console.error('Failed to create keypoint', err);
      },
    });
  }
}
