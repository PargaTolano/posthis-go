import { ReplaySubject, Subject }  from 'rxjs';

const toastSubject = new Subject();

export const toastService = {
    makeToast,
    toast$: toastSubject.asObservable(),
};

function makeToast( content, type ){
    toastSubject.next( {content, type });
}