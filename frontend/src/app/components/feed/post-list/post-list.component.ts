// /frontend/src/app/components/feed/post-list/post-list.component.ts
import { Component, Input } from '@angular/core';
import { Post } from '../../../models/post.model';

@Component({
  selector: 'app-post-list',
  template: `
    <div *ngFor="let post of posts" class="post-card">
      <mat-card>
        <mat-card-header>
          <mat-card-title>{{ post.username }}</mat-card-title>
          <mat-card-subtitle>
            {{ post.created_at | date: 'short' }}
          </mat-card-subtitle>
        </mat-card-header>
        <mat-card-content>
          <p>{{ post.content }}</p>
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [
    `
      .post-card {
        margin-bottom: 15px;
      }
    `,
  ],
})
export class PostListComponent {
  @Input() posts: Post[] = [];
}
