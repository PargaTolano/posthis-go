import React from 'react';

import { CircularProgress } from '@material-ui/core';

import { loadingState, useLoadingState }          from '_hooks';

import styles from '_styles/Loading.module.css';

export const Loading = () => {

    useLoadingState();

    if (!loadingState.get)
        return null;

    return (
        <div className={styles.loading}>
            <CircularProgress color='primary'/>
        </div>
    );
}

export default Loading;