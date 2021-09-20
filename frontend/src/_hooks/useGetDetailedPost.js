import {useState, useEffect} from 'react';

import { getPost }           from '_api';
import { handleResponse }    from '_helpers';

export const useGetDetailedPost = ( id ) =>{
    
    const [ state, setState ] =  useState([false, null]);

    useEffect(()=>{
        (async()=>{
            const {data:responseData, err} = await getPost( id );

            if ( err !== null){
                setState( [ true, null ] );
                return;
            }

            const {data} = responseData;
            setState([true, data]);
        })();

    }, [id]);
    
    return state;
};

export default useGetDetailedPost;