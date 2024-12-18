import { useState, useEffect, useRef } from "react";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
    Home,
    Settings,
    Users,
    FileText,
    CalendarCheck,
    Package,
    ChevronRight,
} from "lucide-react";
import { cn } from "@/lib/utils";

// Tabs
import Overview from "./dashboard/Overview";
import Schedule from "./schedule/Schedule";
import Resources from "./resources/Resources";
import { Upload, Approval, History } from "./labbooks";
import { AllUsers, Register, Permissions } from "./users";
import {
    GeneralSettings,
    SecuritySettings,
    NotificationsSettings,
    IntegrationsSettings,
} from "./settings";

type NavItem = {
    icon: JSX.Element;
    label: string;
    id: string;
    role: number;
    subTabs: { name: string; role: number }[];
    components: {
        [key: string]: React.ComponentType;
    };
};

type NavigationItemProps = {
    item: NavItem;
    handleItemClick: (itemId: string) => void;
    handleSubTabClick: (itemId: string, subTabName: string) => void;
    currentTab: string;
    currentSubTab: string;
    userRole: number;
};

type NavigationProps = {
    currentTab: string;
    currentSubTab: string;
    userRole: number;
    onItemClick: (itemId: string) => void;
    onSubTabClick: (itemId: string, subTabName: string) => void;
    isDesktop: boolean;
    onMobileClose?: () => void;
};

export const navItems: NavItem[] = [
    {
        icon: <Home className="w-4 h-4" />,
        label: "Dashboard",
        id: "dashboard",
        role: 0,
        subTabs: [
            { name: "Overview", role: 0 },
            // { name: "Analytics", role: 1 },
            // { name: "Reports", role: 2 },
            // { name: "Metrics", role: 2 },
        ],
        components: {
            overview: Overview,
            // analytics: Analytics,
            // reports: Reports,
            // metrics: Metrics,
        },
    },
    {
        icon: <FileText className="w-4 h-4" />,
        label: "Lab Books",
        id: "labbooks",
        role: 1,
        subTabs: [
            { name: "Upload", role: 1 },
            { name: "Approval", role: 2 },
            { name: "History", role: 1 },
        ],
        components: {
            upload: Upload,
            approval: Approval,
            history: History,
        },
    },
    {
        icon: <CalendarCheck className="w-4 h-4" />,
        label: "Schedules",
        id: "schedules",
        role: 1,
        subTabs: [],
        components: {
            default: Schedule,
        },
    },
    {
        icon: <Package className="w-4 h-4" />,
        label: "Resources",
        id: "resources",
        role: 1,
        subTabs: [],
        components: {
            default: Resources,
        },
    },
    {
        icon: <Users className="w-4 h-4" />,
        label: "Users",
        id: "users",
        role: 2,
        subTabs: [
            { name: "All Users", role: 2 },
            { name: "Register", role: 2 },
            { name: "Permissions", role: 3 },
        ],
        components: {
            "all users": AllUsers,
            register: Register,
            permissions: Permissions,
        },
    },
    {
        icon: <Settings className="w-4 h-4" />,
        label: "Settings",
        id: "settings",
        role: 1,
        subTabs: [
            { name: "General", role: 1 },
            { name: "Security", role: 1 },
            { name: "Notifications", role: 3 },
            { name: "Integrations", role: 3 },
        ],
        components: {
            general: GeneralSettings,
            security: SecuritySettings,
            notifications: NotificationsSettings,
            integrations: IntegrationsSettings,
        },
    },
];

const NavigationItem: React.FC<NavigationItemProps> = ({
    item,
    handleItemClick,
    handleSubTabClick,
    currentTab,
    currentSubTab,
    userRole,
}) => {
    const [isHovered, setIsHovered] = useState(false);
    const [dropdownTop, setDropdownTop] = useState<number | null>(null);
    const itemRef = useRef<HTMLDivElement>(null);
    const timeoutRef = useRef<NodeJS.Timeout>();

    const hasSubTabs = item.subTabs.some((subTab) => userRole >= subTab.role);
    const filteredSubTabs = item.subTabs.filter(
        (subTab) => userRole >= subTab.role
    );

    const handleClick = (e: React.MouseEvent) => {
        if (hasSubTabs) {
            e.preventDefault();
            e.stopPropagation();
            return;
        }
        handleItemClick(item.id);
    };

    const updateDropdownPosition = () => {
        if (itemRef.current) {
            const rect = itemRef.current.getBoundingClientRect();
            setDropdownTop(rect.top);
        }
    };

    useEffect(() => {
        if (itemRef.current) {
            updateDropdownPosition();
        }

        const handleScroll = () => {
            if (isHovered) {
                updateDropdownPosition();
            }
        };

        window.addEventListener("scroll", handleScroll, true);
        return () => window.removeEventListener("scroll", handleScroll, true);
    }, [isHovered]);

    const handleMouseEnter = () => {
        if (timeoutRef.current) {
            clearTimeout(timeoutRef.current);
        }
        updateDropdownPosition();
        setIsHovered(true);
    };

    const handleMouseLeave = (e: React.MouseEvent) => {
        const relatedTarget = e.relatedTarget as HTMLElement;
        const isMovingToChild = itemRef.current?.contains(relatedTarget);

        if (!isMovingToChild) {
            timeoutRef.current = setTimeout(() => {
                setIsHovered(false);
            }, 50);
        }
    };

    return (
        <div
            ref={itemRef}
            className="relative"
            onMouseEnter={handleMouseEnter}
            onMouseLeave={handleMouseLeave}
        >
            <div className="px-2 py-1">
                <div
                    onClick={handleClick}
                    className={cn(
                        "w-full flex items-center rounded-lg px-3 py-2 transition-all duration-200 select-none",
                        currentTab === item.id || isHovered
                            ? "bg-gray-100 text-gray-900"
                            : "text-gray-500 hover:text-gray-900 hover:bg-gray-100",
                        hasSubTabs ? "cursor-default" : "cursor-pointer"
                    )}
                >
                    <div className="flex items-center gap-3 flex-1">
                        {item.icon}
                        <span>{item.label}</span>
                    </div>
                    {hasSubTabs && (
                        <ChevronRight className="h-4 w-4 ml-2 flex-shrink-0" />
                    )}
                </div>
            </div>

            {hasSubTabs && isHovered && dropdownTop !== null && (
                <>
                    <div
                        className="fixed h-full"
                        style={{
                            left: "15rem",
                            top: dropdownTop,
                            width: "2rem",
                            zIndex: 49,
                        }}
                        onMouseEnter={handleMouseEnter}
                        onMouseLeave={handleMouseLeave}
                    />
                    <div
                        className={cn(
                            "fixed w-48 bg-white border rounded-lg shadow-lg",
                            "animate-in zoom-in-95 duration-200",
                            "origin-top-left"
                        )}
                        style={{
                            zIndex: 50,
                            left: "16rem",
                            top: dropdownTop,
                            transform: `translateZ(0)`,
                        }}
                        onMouseEnter={handleMouseEnter}
                        onMouseLeave={handleMouseLeave}
                    >
                        {filteredSubTabs.map((subTab) => (
                            <div
                                key={subTab.name}
                                onClick={() => {
                                    handleSubTabClick(item.id, subTab.name);
                                    setIsHovered(false);
                                }}
                                className={cn(
                                    "px-3 py-2 text-sm hover:bg-gray-100 cursor-pointer first:rounded-t-lg last:rounded-b-lg",
                                    currentTab === item.id &&
                                        currentSubTab === subTab.name
                                        ? "bg-gray-100"
                                        : ""
                                )}
                            >
                                {subTab.name}
                            </div>
                        ))}
                    </div>
                </>
            )}
        </div>
    );
};

export const Navigation: React.FC<NavigationProps> = ({
    currentTab,
    currentSubTab,
    userRole,
    onItemClick,
    onSubTabClick,
    isDesktop,
    onMobileClose,
}) => {
    const handleItemClick = (itemId: string) => {
        onItemClick(itemId);
        if (!isDesktop && onMobileClose) {
            onMobileClose();
        }
    };

    const handleSubTabClick = (itemId: string, subTabName: string) => {
        onSubTabClick(itemId, subTabName);
        if (!isDesktop && onMobileClose) {
            onMobileClose();
        }
    };

    return (
        <div className="h-full flex flex-col relative z-50">
            <div className="h-14 flex flex-col items-center justify-center px-4 border-b">
                <h2 className="text-lg font-bold">AlphaLab</h2>
                <span className="text-xs text-gray-500">v0.0.1</span>
            </div>
            <ScrollArea className="flex-1 py-2">
                <nav className="space-y-1 px-2">
                    {navItems
                        .filter((item) => userRole >= item.role)
                        .map((item) => (
                            <NavigationItem
                                key={item.id}
                                item={item}
                                handleItemClick={handleItemClick}
                                handleSubTabClick={handleSubTabClick}
                                currentTab={currentTab}
                                currentSubTab={currentSubTab}
                                userRole={userRole}
                            />
                        ))}
                </nav>
            </ScrollArea>
        </div>
    );
};

export const getComponent = (
    tabId: string,
    subTab?: string
): React.ComponentType | null => {
    const navItem = navItems.find((item) => item.id === tabId);
    if (!navItem) return null;

    if (subTab) {
        return navItem.components[subTab.toLowerCase()] || null;
    }

    return navItem.components.default || null;
};
