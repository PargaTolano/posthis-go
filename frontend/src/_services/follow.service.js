import { BehaviorSubject }    from 'rxjs';
import { toastService }     from '_services';

import {
    getFollowers, 
    getFollowing
} from '_api';

const followSubject = new BehaviorSubject({
    open:       false,
    loading:    false,
    title:      'Followers',
    users:      null
});

export const followService = {
    getFollowerUsers,
    getFollowedUsers,
    close,
    follow$: followSubject.asObservable(),
    get currValue(){ return followSubject.value}
};

async function getFollowerUsers( id ) {

    followSubject.next({open: true, loading: true, title: 'Followers', users: null});
    const {data:responseData, err} = await getFollowers(id);

    if ( err !== null ){
        toastService.makeToast('Error on getting followers: ' + err.message, 'error');
        return;
    }

    const { data } = responseData;

    followSubject.next({...followService.currValue, loading: false, users: data});

}

async function getFollowedUsers( id ) {

    followSubject.next({open: true, loading: true, title: 'Following', users: null});
    const {data:responseData, err} = await getFollowing(id);

    if ( err !== null ){
        toastService.makeToast('Error on getting followers: ' + err.message, 'error');
        return;
    }

    const { data } = responseData;

    followSubject.next({...followService.currValue, loading: false, users: data});
}

async function close(event, reason){
    followSubject.next({ ...followService.currValue, open: false, loading: false});
}