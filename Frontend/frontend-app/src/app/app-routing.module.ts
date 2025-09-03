import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { EditUserComponent } from './users/edit-user/edit-user.component';
import { ViewUserComponent } from './users/view-user/view-user.component';
import { AdminPageComponent } from './admin-page/admin-page.component';
import { ForbiddenComponent } from './forbbiden/forbbiden.component';
import { AdminGuard } from './auth/admin.guard';
import { ListBlogsComponent } from './blog/list-blogs/list-blogs.component';
import { CreateBlogComponent } from './blog/create-blog/create-blog.component';
import { CreateTourComponent } from './tour/create-tour/create-tour.component';
import { CreateKeypointComponent } from './tour/create-keypoint/create-keypoint.component';
import { ListToursComponent } from './tour/list-tours/list-tours.component';
import { ViewTourComponent } from './tour/view-tour/view-tour.component';
import { PositionSimComponent } from './tour/position-sim/position-sim.component';

const routes: Routes = [
  { path: '', component: HomeComponent },
  {
    path: 'auth',
    loadChildren: () => import('./auth/auth.module').then((m) => m.AuthModule),
  },
  { path: 'users/:id/edit', component: EditUserComponent },
  { path: 'users/:id/view', component: ViewUserComponent },
  { path: 'admin', component: AdminPageComponent, canActivate: [AdminGuard] },
  { path: 'forbidden', component: ForbiddenComponent },
  { path: 'blogs', component: ListBlogsComponent },
  { path: 'users/:id/create-blog', component: CreateBlogComponent },
  { path: 'users/:id/create-tour', component: CreateTourComponent },
  { path: 'tours/:id/create-keypoint', component: CreateKeypointComponent },
  { path: 'tours', component: ListToursComponent },
  { path: 'users/:id/tours', component: ListToursComponent },
  { path: 'tours/view/:id', component: ViewTourComponent },
  { path: 'position-sim', component: PositionSimComponent},

  { path: '**', redirectTo: '' }, // Redirect any unknown paths to home,
  
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
