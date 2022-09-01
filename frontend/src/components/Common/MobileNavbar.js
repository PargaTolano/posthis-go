import React,{useRef, useState, useEffect}    from 'react';

import { Link } from 'react-router-dom';

import {  
  Search                  as SearchIcon,
  ExitToAppTwoTone        as LogoutIcon,
  Menu                    as MenuIcon
} from '@material-ui/icons';


import profilePicPlaceholder from 'assets/avatar-placeholder.svg';

import { authenticationService } from '_services';

import { routes } from '_utils';

import LogoNominado               from 'assets/logPT.svg';
import Logo from 'assets/Logo.png';

import mobileStyles from '_styles/MobileNavbar.module.css';
import styles from '_styles/NavbarMobile.module.css';

const Sidebar=({history, open, setOpen})=>{
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
            <img className={styles.logo} src= {Logo}/>
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
  );
};

export const MobileNavbar = ({ history }) => {
  const [open, setOpen]=useState(false);
  return (
    <>
      <div 
        className={mobileStyles.darkener}
        onClick={()=>setOpen(x=>!x)}
        data-open={open} 
      >
      </div>
      <nav className={mobileStyles.navbar}>
          <MenuIcon 
            fontSize='inherit' 
            className={mobileStyles.menuIcon}
            onClick={()=>setOpen(true)}
          />
          <Link className={mobileStyles.link} to={routes.feed}>
            <img className={mobileStyles.logo} src= {Logo}/>
          </Link> 
      </nav>
      <Sidebar history={history} open={open} setOpen={setOpen}/>
    </>
  );
};

export default MobileNavbar;