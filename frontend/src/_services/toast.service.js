import { ReplaySubject }  from 'rxjs';

const toastSubject = new ReplaySubject();

export const toastService = {
    makeToast,
    toast$: toastSubject.asObservable(),
};

function makeToast( content, type ){
    toastSubject.next( {content, type });
}