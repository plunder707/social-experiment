// /frontend/src/app/components/feed/feed.component.ts
import { Component, OnInit } from '@angular/core';
import { Post } from '../../models/post.model';
import { WebSocketService } from '../../services/websocket.service';
import { HttpClient } from '@angular/common/http';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-feed',
  template: `
    <app-post-form (newPost)="addPost($event)"></app-post-form>
    <app-post-list [posts]="posts"></app-post-list>
  `,
})
export class FeedComponent implements OnInit {
  posts: Post[] = [];

  constructor(
    private wsService: WebSocketService,
    private http: HttpClient,
    private auth: AuthService
  ) {}

  ngOnInit() {
    this.fetchPosts();

    this.wsService.getNewPosts().subscribe((post: Post) => {
      this.posts.unshift(post);
    });
  }

  fetchPosts() {
    this.http.get<Post[]>('http://localhost:8080/posts').subscribe({
      next: (data) => {
        this.posts = data;
      },
      error: (err) => {
        console.error('Error fetching posts:', err);
      },
    });
  }

  addPost(post: Post) {
    this.posts.unshift(post);
  }
}
