import {
  Component,
  AfterViewInit,
  ViewChild,
  ElementRef,
  OnDestroy,
  OnInit,
} from '@angular/core';
import * as L from 'leaflet';
import 'leaflet-routing-machine';
import { ActivatedRoute, Router } from '@angular/router';
import { KeypointService } from '../keypoint.service';
import { TourService } from '../tour.service';
import { AuthService } from '../../auth/auth.service';
import { MatDialog } from '@angular/material/dialog';
import { KeypointDetailDialogComponent } from '../keypoint-detail-dialog/keypoint-detail-dialog.component';

@Component({
  selector: 'app-view-tour',
  templateUrl: './view-tour.component.html',
  styleUrls: ['./view-tour.component.css'],
})
export class ViewTourComponent implements OnInit, AfterViewInit, OnDestroy {
  private routeControl?: L.Routing.Control;
  @ViewChild('mapContainer', { static: true })
  mapContainer!: ElementRef<HTMLDivElement>;
  private map?: L.Map;
  private markers: L.Marker[] = [];
  private existingIcon: L.Icon = L.icon({
    iconUrl: 'assets/icons/keypointmark.svg',
    iconSize: [32, 32],
    iconAnchor: [16, 32],
  });

  tourId?: string;
  tour: any = null;
  currentUserId?: string;
  currentUserRole?: string;
  isOwner = false;
  isEditMode = false;
  editForm: any = {};
  availableTags = [
    'Nature',
    'History',
    'Adventure',
    'Food',
    'Culture',
    'Relax',
  ];

  constructor(
    private route: ActivatedRoute,
    private keypointService: KeypointService,
    private tourService: TourService,
    private auth: AuthService,
    private dialog: MatDialog,
    private router: Router
  ) {}

  ngOnInit(): void {
    // load current user info first so we can determine ownership
    this.auth.whoAmI().subscribe({
      next: (u) => {
        this.currentUserId = u?.id?.toString();
        this.currentUserRole = u?.role;
        this.initializeFromRoute();
      },
      error: () => {
        this.currentUserId = undefined;
        this.currentUserRole = undefined;
        this.initializeFromRoute();
      },
    });
  }

  private initializeFromRoute() {
    // prefer snapshot for immediate value, also subscribe to changes
    const sid = this.route.snapshot.paramMap.get('id');
    if (sid) {
      this.tourId = sid;
      this.loadTourDetails();
    }

    this.route.paramMap.subscribe((pm) => {
      const id = pm.get('id');
      if (id && id !== this.tourId) {
        this.tourId = id;
        console.log('Tour ID: ' + this.tourId);
        this.loadTourDetails();
        // if map already initialized, reload keypoints
        if (this.map) this.loadKeypoints();
      }
    });
  }

  private loadTourDetails(): void {
    if (!this.tourId) return;
    console.log('Loading tour details for id=' + this.tourId);
    this.tourService.getById(this.tourId).subscribe({
      next: (t) => {
        this.tour = t;
        // determine ownership (authorId may be number or string)
        const authorIdRaw = (t as any)['authorId'] ?? (t as any)['userId'];
        const authorId =
          authorIdRaw != null ? authorIdRaw.toString() : undefined;
        this.isOwner = !!(
          this.currentUserId &&
          authorId &&
          this.currentUserId === authorId
        );
        // if the map is already initialized, load keypoints now
        if (this.map) {
          this.loadKeypoints();
        }
      },
      error: (err) => console.error('Failed to load tour', err),
    });
  }

  ngAfterViewInit(): void {
    this.map = L.map(this.mapContainer.nativeElement, {
      center: [20, 0],
      zoom: 2,
    });
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors',
    }).addTo(this.map);

    setTimeout(() => this.map?.invalidateSize(), 200);
    this.loadKeypoints();
  }

  private loadKeypoints(): void {
    if (!this.tourId || !this.map) return;
    this.keypointService.getByTourSorted(this.tourId).subscribe({
      next: (kps) => {
        this.markers.forEach((m) => m.remove());
        this.markers = [];
        const latlngs: L.LatLng[] = [];
        // If viewer is not owner, only expose the first keypoint
        const kpsToShow =
          !this.isOwner && (kps || []).length > 0
            ? [(kps || [])[0]]
            : kps || [];
        kpsToShow.forEach((kp: any) => {
          const lat = kp.coordinates?.latitude;
          const lng = kp.coordinates?.longitude;
          if (lat == null || lng == null) return;
          latlngs.push(L.latLng(lat, lng));
          const m = L.marker([lat, lng], { icon: this.existingIcon }).addTo(
            this.map!
          );
          m.on('click', () => {
            const ref = this.dialog.open(KeypointDetailDialogComponent, {
              data: kp,
              width: '320px',
            });
            ref.afterClosed().subscribe((res) => {
              if (!res) return;
              if (res.action === 'deleted' || res.action === 'delete') {
                this.loadKeypoints();
                this.updateTourDistanceAndDuration();
              }
            });
          });
          this.markers.push(m);
        });
        // adjust map view to markers if any
        if (this.markers.length > 0) {
          const group = L.featureGroup(this.markers);
          try {
            this.map!.fitBounds(group.getBounds().pad(0.3));
          } catch (e) {
            // ignore fitBounds errors
            console.warn('fitBounds failed', e);
          }
        }
        // Draw street-connected route using Leaflet Routing Machine
        if (latlngs.length > 1) {
          if (this.routeControl) {
            this.map!.removeControl(this.routeControl);
            this.routeControl = undefined;
          }
          this.routeControl = L.Routing.control({
            waypoints: latlngs,
            router: L.routing.mapbox(
              'pk.eyJ1IjoidmVsam9vMDIiLCJhIjoiY20yaGV5OHU4MDFvZjJrc2Q4aGFzMTduNyJ9.vSQUDO5R83hcw1hj70C-RA',
              { profile: 'mapbox/walking' }
            ),
            show: false,
            lineOptions: {
              styles: [{ color: 'blue', weight: 4 }],
              extendToWaypoints: false,
              missingRouteTolerance: 0,
            },
            createMarker: () => null,
          } as any).addTo(this.map!);
        }
      },
      error: (err) => console.error('failed to load keypoints', err),
    });
  }

  ngOnDestroy(): void {
    if (this.routeControl && this.map) {
      this.map.removeControl(this.routeControl);
      this.routeControl = undefined;
    }
    if (this.map) {
      this.map.remove();
      this.map = undefined;
    }
  }

  difficultyLabel(d: number): string {
    // backend: 0=Beginner,1=Intermediate,2=Advanced,3=Pro
    switch (d) {
      case 0:
        return 'Beginner';
      case 1:
        return 'Intermediate';
      case 2:
        return 'Advanced';
      case 3:
        return 'Pro';
      default:
        return 'Unknown';
    }
  }

  goToCreateKeypoint() {
    // navigate to route that allows creating keypoint for a tour
    this.router.navigate(['/tours', this.tourId, 'create-keypoint']);
  }

  goToMyTours() {
    this.auth.whoAmI().subscribe({
      next: (me) => this.router.navigate(['/users', me.id, 'tours']),
      error: () => this.router.navigate(['/auth/login']),
    });
  }

  enableEditMode() {
    this.isEditMode = true;
    // Initialize edit form with current tour data
    this.editForm = {
      title: this.tour.title,
      description: this.tour.description,
      difficulty: this.difficultyLabel(this.tour.difficulty),
      selectedTags: this.availableTags.map((tag) =>
        this.tour.tags.includes(tag)
      ),
      price: this.tour.price,
    };
  }

  cancelEdit() {
    this.isEditMode = false;
    this.editForm = {};
  }

  saveChanges() {
    if (!this.tourId || !this.tour) return;

    // Build selected tags array from checkboxes
    const selectedTags = this.editForm.selectedTags
      .map((selected: boolean, index: number) =>
        selected ? this.availableTags[index] : null
      )
      .filter((tag: string | null) => tag !== null);

    const updatedTourDto = {
      id: this.tour.id,
      authorId: this.tour.authorId,
      title: this.editForm.title,
      description: this.editForm.description,
      tags: selectedTags,
      price: parseFloat(this.editForm.price),
      distance: this.tour.distance,
      duration: this.tour.duration,
      status: this.statusToNumber(this.tour.status),
      difficulty: this.difficultyToNumber(this.editForm.difficulty),
      transportType: this.transportTypeToNumber(this.tour.transportType),
      publishedAt: this.tour.publishedAt,
      archivedAt: this.tour.archivedAt,
    };

    this.tourService.update(this.tourId, updatedTourDto).subscribe({
      next: () => {
        // Update local tour object
        this.tour.title = this.editForm.title;
        this.tour.description = this.editForm.description;
        this.tour.difficulty = this.editForm.difficulty;
        this.tour.tags = selectedTags;
        this.tour.price = updatedTourDto.price;

        this.isEditMode = false;
        this.editForm = {};
        console.log('Tour updated successfully');
      },
      error: (err) => {
        console.error('Failed to update tour', err);
        alert('Failed to update tour. Please try again.');
      },
    });
  }

  // Helper functions to convert string enums to numbers for backend
  private statusToNumber(status: string): number {
    switch (status) {
      case 'Draft':
        return 0;
      case 'Published':
        return 1;
      case 'Archived':
        return 2;
      default:
        return 0;
    }
  }

  private difficultyToNumber(difficulty: string): number {
    switch (difficulty) {
      case 'Beginner':
        return 0;
      case 'Intermediate':
        return 1;
      case 'Advanced':
        return 2;
      case 'Pro':
        return 3;
      default:
        return 0;
    }
  }

  private transportTypeToNumber(transportType: string): number {
    switch (transportType) {
      case 'Walking':
        return 0;
      case 'Bicycle':
        return 1;
      case 'Bus':
        return 2;
      default:
        return 0;
    }
  }

  // Helper to calculate distance between two lat/lng points (Haversine formula)
  private haversineDistance(
    lat1: number,
    lng1: number,
    lat2: number,
    lng2: number
  ): number {
    const toRad = (x: number) => (x * Math.PI) / 180;
    const R = 6371; // Earth radius in km
    const dLat = toRad(lat2 - lat1);
    const dLng = toRad(lng2 - lng1);
    const a =
      Math.sin(dLat / 2) * Math.sin(dLat / 2) +
      Math.cos(toRad(lat1)) *
        Math.cos(toRad(lat2)) *
        Math.sin(dLng / 2) *
        Math.sin(dLng / 2);
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    return R * c;
  }

  // Update tour's distance and duration after keypoint deletion
  private updateTourDistanceAndDuration(): void {
    if (!this.tourId) return;
    this.keypointService.getByTourSorted(this.tourId).subscribe({
      next: (keypoints) => {
        let totalDistance = 0;
        for (let i = 1; i < keypoints.length; i++) {
          const prev = keypoints[i - 1].coordinates;
          const curr = keypoints[i].coordinates;
          if (prev && curr) {
            totalDistance += this.haversineDistance(
              prev.latitude,
              prev.longitude,
              curr.latitude,
              curr.longitude
            );
          }
        }
        // Example: duration = totalDistance * 12 (assume 12 min per km, adjust as needed)
        const totalDuration = Math.round(totalDistance * 12);
        // Update tour
        if (this.tour && this.tourId) {
          const updatedTourDto = {
            id: this.tour.id,
            authorId: this.tour.authorId,
            title: this.tour.title,
            description: this.tour.description,
            tags: this.tour.tags,
            price: this.tour.price,
            distance: Math.round(totalDistance * 100) / 100, // round to 2 decimals
            duration: totalDuration,
            status: this.statusToNumber(this.tour.status),
            difficulty: this.difficultyToNumber(this.tour.difficulty),
            transportType: this.transportTypeToNumber(this.tour.transportType),
            publishedAt: this.tour.publishedAt,
            archivedAt: this.tour.archivedAt,
          };
          this.tourService.update(this.tourId, updatedTourDto).subscribe({
            next: () => {
              this.tour.distance = updatedTourDto.distance;
              this.tour.duration = updatedTourDto.duration;
              console.log('Tour distance and duration updated successfully');
            },
            error: (err) =>
              console.error('Failed to update tour distance/duration', err),
          });
        }
      },
      error: (err) =>
        console.error('Failed to recalculate tour distance/duration', err),
    });
  }
}
