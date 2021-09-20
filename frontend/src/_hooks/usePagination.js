import { useRef, useState, useEffect }  from 'react';

export const usePagination = ( hasFetched, total, last, limit, onIntersection, rootMargin = '0px') =>{

    const [ more, setMore ] = useState( true );
        
    const ref = useRef();  

    useEffect( () => {
        if ( last < limit ){
            setMore(!hasFetched);
        }
    }, [ hasFetched, total, last, limit ]);

    useEffect( () => {
        let observer = new IntersectionObserver(
            (entries, obs) => {

                entries.forEach( entry =>{
                    if (entry.isIntersecting) onIntersection(); 
                });
            },
            {
                root: null,
                rootMargin,
            });

        observer.observe(ref.current);


        return () => void observer.disconnect();
    },[]);

    return {ref, more};
};

export default usePagination;