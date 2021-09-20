import { BehaviorSubject }  from 'rxjs';
import { authTokenKey }     from '_utils';
import { toastService } from '_services'

import { logIn }            from '_api/users';

const currentUserSubject = new BehaviorSubject( JSON.parse( localStorage.getItem( authTokenKey ) ) );

export const authenticationService = {
    login,
    logout,
    currentUser: currentUserSubject.asObservable(),
    get currentUserValue () { return currentUserSubject.value }
};

function login( username, password ) {

    let data = { username, password};

    return logIn(data)
        .then(({data, err}) => {

            if( err ){
                toastService.makeToast( err.message, 'error');
                return null;
            }

            let {data:userData} = data;
            localStorage.setItem( authTokenKey, JSON.stringify( userData ) );
            currentUserSubject.next( userData );
            
            return data;
        });
}

function logout() {
    localStorage.removeItem( authTokenKey );
    currentUserSubject.next(null);
}