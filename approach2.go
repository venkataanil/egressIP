func find_index(new_ip int, sorted_excluded_ips []int) {
    return index if new_ip present in sorted_excluded_ips
    return -1 if doesn’t match
} 

# sort given exclude_ips
var sorted_excluded_ips []int
for i := range exclude_ips {
   sorted_excluded_ips = append(excluded_ip_list, convert_into_int(i))
}

cidr_start_int := convert_cidr_to_int(cidr)
excluded_ip_index := 0  

# list of ips for this current object in the iteration
var ips_for_obj []string

for i = range addr_per_iteration {
   // we have already used ((iteration * addr_per_iteration) + excluded_ip_index) in previous iterations
   new_ip = cidr_start_int + (iteration * addr_per_iteration)
   index := find_index(new_ip, sorted_excluded_ips)
   If index != -1 {
      excluded_ip_index = index
   }
   new_ip := new_ip + excluded_ip_index
   ips_for_obj = append( Ips_for_obj, convert_ip_int_tostring(new_ip)
}
