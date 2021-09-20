import React,{ useState, useEffect} from 'react';

import {DialogFollow, Followers} from 'components/Follow';

import { followService } from '_services';

export const DialogFollowWrapper = ({history}) => {

    const [{open, loading, title, users}, setState] = 
    useState({
      open:       false,
      loading:    false,
      title:      'Followers',
      users:      null
    });

    useEffect(()=>{
        const subs = 
        followService.follow$.subscribe( x => void setState(x));

        return ()=> subs.unsubscribe();
    },[]);


    return (
        <DialogFollow
          open={open}
          title={title}
          onClose={followService.close}
        >
          <Followers loading={loading} users={users}/>
        </DialogFollow>
    )
}

export default DialogFollowWrapper
