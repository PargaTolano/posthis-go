
//////////////////////////////////
//                              //
//            UNUSED            //
//                              //
//////////////////////////////////


import { BehaviorSubject }  from 'rxjs';
import { handleResponse }   from '_helpers';

import { getPost ,getFeed, getUserFeed }            from '_api/posts';

const postSubject = new BehaviorSubject( null );

export const postService = {
    getDetailPost,
    getFeedPosts,
    getUserFeedPosts,
    post$: postSubject.asObservable(),
};

function getDetailPost(id){
    return getPost(id)
        .then(handleResponse)
        .then(res=>{
            const { data } = res;
            postSubject.next( data );

            return data;
        })
}

function getFeedPosts( offset, limit ) {

    return getFeed(offset, limit)
        .then(handleResponse)
        .then(res => {
            const { data } = res;
            postSubject.next( data );
            
            return data;
        });
}

function getUserFeedPosts( id, offset, limit ) {

    return getUserFeed( id, offset, limit)
    .then(handleResponse)
    .then(res => {
        const { data } = res;
        postSubject.next( data );
        
        return data;
    });
}