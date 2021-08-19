import {useEffect, useState} from 'react';

export let loadingState = {get: null, set: null}

export const useLoadingState = ()=>{
    const [get, set] = useState(false);

    loadingState = { get, set};
}