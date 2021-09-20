import React                    from 'react';
import { Redirect }             from 'react-router-dom';

import { NavBar }               from 'components/Feed';
import { ProfileContainer }     from 'components/Profile';
import { DialogFollowWrapper } from 'components/Follow';

import { routes }               from '_utils';
import { useGetUserProfile }    from '_hooks';

import styles from '_styles/ProfileDetail.module.css';
export const ProfileDetail = ({ match, history, ...rest }) => {
  
  const { id }  = match.params;

  const [[ready, user], setUser] = useGetUserProfile( id || '' );

  if( id == 'undefined' || id === undefined || id === null || id === ''|| (ready && user === null ) ){
    return (<Redirect to={routes.feed}/>);
  }

  return (
      <div className = { styles.root }>

        <DialogFollowWrapper history={history}/>
        <NavBar history={history} {...rest}/>
        {
          ready && <ProfileContainer user={user} setUser={setUser} {...rest}/>
        }
      </div>
  );
};

export default ProfileDetail;
