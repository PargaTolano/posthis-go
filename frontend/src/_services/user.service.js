import { BehaviorSubject }  from 'rxjs';
import { handleResponse }   from '_helpers';

const userSubject = new BehaviorSubject( {} );

export const userService = {
    user$: userSubject.asObservable(),
};