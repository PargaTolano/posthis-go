import React from 'react';

import { CircularProgress } from '@material-ui/core';

import { usePagination }          from '_hooks';

import styles from '_styles/PaginationElement.module.css';

const defaultMessage = 'Pagination Element Scrolled into view';
const defaultFunc = () => void console.log(defaultMessage);

export const PaginationElement = ({
    name , 
    hasFetched, 
    total, 
    last,
    limit,
    onIntersection = defaultFunc,
    rootMargin     = '0px'}) => {

    const {ref, more} = usePagination(hasFetched, total, last, limit, onIntersection, rootMargin);

    return (
        <div 
            ref={ref}
            className={styles.loading}
        >

            {
                more ?
                    <CircularProgress color='primary'/>

                :
                    <h3 className={styles.unavailable}>Sorry, no more {name} are available right now</h3>

            }
        </div>
    );
}

export default PaginationElement;