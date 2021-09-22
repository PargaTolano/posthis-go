import React                from 'react';

import { useMediaQuery }    from 'react-responsive';

import { MobileNavbar }     from 'components/Common/MobileNavbar';
import { NavBar }           from 'components/Common/Navbar';

export const NavbarWrapper = (props) => {
    const isMobile =  useMediaQuery({ query: '(max-width: 768px)' });
    return (
        <>
        {
            isMobile ?
            <MobileNavbar {...props}/>
            :
            <NavBar {...props}/>
        }
        </>
    )
}

export default NavbarWrapper;
