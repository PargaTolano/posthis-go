import React,{useRef, useState, useEffect}    from 'react';

import { Link } from 'react-router-dom';

import {  
    Search                  as SearchIcon,
    ExitToAppTwoTone        as LogoutIcon
} from '@material-ui/icons';

import styles from '_styles/NavbarMobile.module.css';

import profilePicPlaceholder from 'assets/avatar-placeholder.svg';

import { authenticationService } from '_services';

import { routes } from '_utils';

import LogoNominado               from 'assets/logPT.svg';

export const MobileNavbar = ({ history, open = true}) => {

    const ref = useRef();

    const [query, setQuery] = useState('');
    
    const onChange = (e)=>void setQuery(e.target.value);

    const onSearch = ()=>{
    };

    const onClickLogOut = ()=>{
      if(window.confirm('Are you sure you want to logout?')){
        authenticationService.logout();
        history.push( routes.login );
      }
    };

    return (
        <div className={ open ? `${styles.root} ${styles.open}` : styles.root}>
            <section className={styles.section}>
              <Link to={routes.feed}>
                <img className={styles.logo} src= {LogoNominado}/>
              </Link>
              <Link
                className={styles.pfpContainer} 
                to={routes.getProfile( authenticationService.currentUserValue.id ) }
              >
                <img 
                  className={styles.pfp} 
                  src={ authenticationService.currentUserValue.profilePicPath || profilePicPlaceholder}
                />
              </Link>
            </section>
            <section className={styles.section}>
              <div className={styles.searchBar}>
                <div className={styles.searchIcon}>
                  <SearchIcon />
                </div>
                <form className={styles.searchBarForm} onSubmit ={onSearch}>
                  <input
                    className={styles.searchBarInput}
                    ref={ref}
                    placeholder='Search…'
                    value={query}
                    onChange  ={onChange}
                  />
                </form>
              </div>
            </section>
            <section className={styles.section}> 
                <button 
                  className={styles.logOutBtn} 
                    onClick={onClickLogOut}
                  >
                    Log out 

                    <LogoutIcon className={styles.logOutIcon}/>
                  </button>
            </section>

        </div>
    )
}

export default MobileNavbar
