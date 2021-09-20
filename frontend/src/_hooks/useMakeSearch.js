import { useState, useEffect }  from 'react';
import { getSearch }            from '_api';
import { handleResponse }       from '_helpers';

import { loadingState }         from '_hooks';
import { toastService } from '_services';

export const useMakeSearch = ( query ) =>{

    const [ state, setState ] = 
        useState([
            false,
            null
        ]);
        
        useEffect( () => {
            
            (async () =>{

                loadingState.set(true);
                const {data:responseData, err} = await getSearch( query, true, true, 0, 5, 0, 5);
                loadingState.set(false);

                if ( err !== null ) {
                    toastService
                        .makeToast( ["Error retrieving search results",err].join('/n'), "error");
                    setState( [true, err] );
                    return;
                }

                const { data } = responseData;

                setState( [true, data] );
            })();

        }, [ query ]);

    return state;
};

export default useMakeSearch;