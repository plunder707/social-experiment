// /frontend/src/app/components/feed/post-form/post-form.component.ts
import { Component, EventEmitter, Output } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { Post } from '../../../models/post.model';
import { AuthService } from '../../../services/auth.service';

@Component({
  selector: 'app-post-form',
  template: `
    <mat-card>
      <form [formGroup]="postForm" (ngSubmit)="onSubmit()">
        <mat-form-field class="full-width">
          <textarea
            matInput
            placeholder="What's on your mind?"
            formControlName="content"
            rows="3"
          ></textarea>
          <mat-error *ngIf="postForm.get('content')?.invalid">
            Post content cannot be empty
          </mat-error>
        </mat-form-field>
        <button
          mat-raised-button
          color="primary"
          type="submit"
          [disabled]="postForm.invalid"
        >
          Post
        </button>
      </form>
    </mat-card>
  `,
  styles: [
    `
      mat-card {
        margin-bottom: 20px;
      }

      .full-width {
        width: 100%;
      }
    `,
  ],
})
export class PostFormComponent {
  @Output() newPost = new EventEmitter<Post>();
  postForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private http: HttpClient,
    private auth: AuthService
  ) {
    this.postForm = this.fb.group({
      content: ['', Validators.required],
    });
  }

  onSubmit() {
    if (this.postForm.invalid) return;

    const { content } = this.postForm.value;

    this.http
      .post<Post>('http://localhost:8080/posts', { content })
      .subscribe({
        next: (post) => {
          this.newPost.emit(post);
          this.postForm.reset();
        },
        error: (err) => {
          console.error('Error creating post:', err);
        },
      });
  }
}
