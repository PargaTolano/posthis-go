import React, { useEffect } from 'react';

import { routes }        from '_utils';

import styles from '_styles/NotFound.module.css';

export const NotFound = ( props ) => {

    const {history} = props;

    useEffect(()=>{
        setTimeout(()=>history.replace(routes.feed), 3000);
    },[]);

    return (
        <div className={styles.root}>
            <div className={styles.message}>
                <h1 className={styles.error}>            ERROR  <span className={styles.errorCode}>404</span>       </h1>
                <p className={styles.errorDescription}>  Page Not Found                                             </p>
                <p className={styles.errorAction}>       Redirecting to home page...                                </p>
            </div>
        </div>
    );
}

export default NotFound;