import { BehaviorSubject }  from 'rxjs';
import { handleResponse }   from '_helpers';

import { getReplies }       from '_api/replies';

const replySubject = new BehaviorSubject( null );

export const replyService = {
    getPostReplies,
    reply$: replySubject.asObservable(),
};

async function getPostReplies(id){

    const {data:responseData, err} = await getReplies( id );

    if ( err !== null ) return;

    const {data} = responseData;

    if ( data !== null )
    replySubject.next(data);
    
    return data;
}