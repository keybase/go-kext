#ifndef KEXT_MANAGER_SAFE_H
#define KEXT_MANAGER_SAFE_H

#include "CoreFoundationSafeTypes.h"

#include <IOKit/kext/KextManager.h>

CFURLRef CFURLCreateWithFileSystemPathSafe(CFAllocatorRef allocator, CFStringRefSafe filePath, CFURLPathStyle pathStyle, Boolean isDirectory) {
  return CFURLCreateWithFileSystemPath(allocator, (CFStringRef)filePath, pathStyle, isDirectory);
}

OSReturn KextManagerLoadKextWithIdentifierSafe(CFStringRefSafe kextIdentifier, CFArrayRefSafe dependencyKextAndFolderURLs) {
  return KextManagerLoadKextWithIdentifier((CFStringRef)kextIdentifier, (CFArrayRef)dependencyKextAndFolderURLs);
}

OSReturn KextManagerUnloadKextWithIdentifierSafe(CFStringRefSafe kextIdentifier) {
  return KextManagerUnloadKextWithIdentifier((CFStringRef)kextIdentifier);
}

CFDictionaryRefSafe KextManagerCopyLoadedKextInfoSafe(CFArrayRefSafe kextIdentifiers, CFArrayRefSafe infoKeys) {
  return (CFDictionaryRefSafe)KextManagerCopyLoadedKextInfo((CFArrayRef)kextIdentifiers, (CFArrayRef)infoKeys);
}

#endif
