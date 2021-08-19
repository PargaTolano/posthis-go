import React from 'react';

import {
    CssBaseline,
    Container,
    CircularProgress
}from  '@material-ui/core';

import { FollowCard } from 'components/Follow';

import styles from '_styles/Follow.module.css';

export const Followers = ( {loading, users} ) => {
    
    return (
        <div 
            className = { styles.root }
        >
            <CssBaseline />
            
            <div className={ loading ? `${styles.paper} ${styles.loading}` : styles.paper}>
                {
                    loading 
                    ?
                        <CircularProgress color='primary'/>
                    :
                        users?.map((v,i)=><FollowCard key={v.id} user={v}/>)
                }
            </div>    
        </div>
    );
  }
  
  export default  Followers;