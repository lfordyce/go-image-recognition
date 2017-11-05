import { Injectable } from '@angular/core';
import { HttpClient, HttpRequest, HttpEvent } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';

// RxJs operator:
import 'rxjs/add/operator/map';

// Model
import { Post } from '../model/post';

@Injectable()
export class FileUploadService {

  private serviceUrl = 'http://localhost:8080';

  constructor(private httpClient: HttpClient) {}

  getApi(): Observable<string> {
    return this.httpClient.get<string>(`${this.serviceUrl}/api`);
  }

  uploadFile(file: File): Observable<HttpEvent<{}>> {
    const formdata: FormData = new FormData();

    formdata.append('image', file);

    const req = new HttpRequest('POST', `${this.serviceUrl}/recognize`, formdata, {
      reportProgress: true,
      responseType: 'json'
    });

    return this.httpClient.request(req);
  }

  analyzePhoto(file: File): Observable<Post> {
    const formdata: FormData = new FormData();

    formdata.append('image', file);

    return this.httpClient.post<Post>(`${this.serviceUrl}/recognize`, formdata);
  }
}
