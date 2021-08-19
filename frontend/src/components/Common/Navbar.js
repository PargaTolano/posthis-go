import React,{ useRef, useState, useEffect} from 'react';

import { Link }                   from 'react-router-dom';

import {  
  Search as SearchIcon, 
  ArrowDropDown as ArrowDropDownIcon 
} from '@material-ui/icons';

import Logo                       from 'assets/Logo.png';
import LogoNominado               from 'assets/logPT.svg';

import { authTokenKey, routes }   from '_utils';
import { authenticationService }  from '_services';

import profilePicPlaceholder from 'assets/avatar-placeholder.svg';

import styles from '_styles/Navbar.module.css';

export const DropDownMenu = ( { visible, setVisible} ) =>{

  const ref = useRef(null);

  useEffect(()=>{

    const clickListener = e=> {
      if( e.target != ref.current){
        setVisible(false);
        return;
      }
    };

    window.addEventListener('click', clickListener);

    return ()=>{
      window.removeEventListener('click', clickListener)
    };
  },[]);

  return (
    <div 
      ref = {ref}
      className={ visible ? styles.dropDownContainer : styles.dropDownContainerInvisible} 
      onClick={ e => void e.stopPropagation()}
    >
      <ul className={styles.dropDownMenu}>

        <li 
          className={styles.dropDownMenuItem}
        >
          <Link
            className={styles.dropDownMenuItemContent}
            to={routes.getProfile(authenticationService.currentUserValue.id)}
          >
            Profile
          </Link>
        </li>

        <li 
          className={styles.dropDownMenuItem}
        >
          <span
            className={styles.dropDownMenuItemContent}
            onClick={authenticationService.logout}
          >
            Log out
          </span>
        </li>
      </ul>
    </div>
  )
};

export const NavBar = ( props ) => {

  const { history } = props;

  const {profilePicPath} = JSON.parse(localStorage.getItem( authTokenKey ));

  const ref = useRef(null);
  const [dropDownVisible, setDropDownVisible] = useState(false);
  const [query, setQuery]       = useState('');

  const handleClick = (e) => { e.stopPropagation(); setDropDownVisible(true) };

  const onChange = ( e ) => void setQuery( e.target.value);

  const onSearch = ( e ) => {
    e.preventDefault();
    ref.current.blur();
    history.push( routes.getSearch(query));
  };

  return (
      <div className={styles.navbar}>

        <div className={styles.content}>
          <Link to={routes.feed}>
            <img className={styles.logo} src= {Logo}/>
            <img className={styles.logo} src= {LogoNominado}/>
          </Link>

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

          <div className={styles.pfpFlex}>
            <Link
              className={styles.pfpContainer} 
              to={routes.getProfile( authenticationService.currentUserValue.id ) }
            >
              <img 
                className={styles.pfp} 
                src={profilePicPath || profilePicPlaceholder}
              />
            </Link>

            <div 
              className={styles.dropDownIconContainer}
              onClick={handleClick}
            >
              <ArrowDropDownIcon 
                className={styles.dropDownIcon}
              />
            </div>
            
            <DropDownMenu
              visible={dropDownVisible}
              setVisible={setDropDownVisible}
            />
          </div>
        </div>
            
      </div>
  );
}

export default NavBar;