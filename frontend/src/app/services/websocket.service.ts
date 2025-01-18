// /frontend/src/app/services/websocket.service.ts
import { Injectable } from '@angular/core';
import { Observable, Subject, webSocket, WebSocketSubject } from 'rxjs';
import { Post } from '../models/post.model';

@Injectable({
  providedIn: 'root',
})
export class WebSocketService {
  private socket$: WebSocketSubject<any>;
  private postSubject = new Subject<Post>();

  constructor() {
    this.connect();
  }

  private connect() {
    this.socket$ = webSocket('ws://localhost:8080/ws');

    this.socket$.subscribe(
      (msg) => {
        const post: Post = msg;
        this.postSubject.next(post);
      },
      (err) => {
        console.error('WebSocket error:', err);
        setTimeout(() => {
          this.connect();
        }, 3000); // Reconnect after 3 seconds
      },
      () => {
        console.warn('WebSocket connection closed');
        setTimeout(() => {
          this.connect();
        }, 3000); // Reconnect after 3 seconds
      }
    );
  }

  getNewPosts(): Observable<Post> {
    return this.postSubject.asObservable();
  }

  sendMessage(msg: any) {
    this.socket$.next(msg);
  }

  close() {
    this.socket$.complete();
  }
}
