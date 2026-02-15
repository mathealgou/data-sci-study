import random
import math

class Client:
	def __init__(self):
		self.age = random.randint(20, 80)
		self.income = random.randint(0, 100_000)
		self.education = random.choice(['high school', 'bachelor', 'master'])
	def __str__(self):
		return f'Client(age={self.age}, income={self.income}, education={self.education}, has_canceled={self.has_canceled})'
	def _get_has_canceled(self):
		# let's say lower education increases the likelihood of cancellation
		likelihood = 0.1  # base likelihood
		if self.education == 'high school':
			likelihood += 0.2
		elif self.education == 'bachelor':
			likelihood += 0.1
		elif self.education == 'master':
			likelihood += 0.05
		# let's say lower income increases the likelihood of cancellation
		likelihood += (50_000 - self.income) / 100_000 * 0.3  # up to 0.3 increase for low income
		# let's say younger age increases the likelihood of cancellation
		likelihood += (60 - self.age) / 40 * 0.2  # up to 0.2 increase for younger age
		has_canceled = random.random() < likelihood
		return has_canceled
	has_canceled = property(_get_has_canceled)

def calculate_entropy_for_attribute(clients, attribute):
	# count the occurrences of each value for the attribute
	value_counts = {}
	for client in clients:
		value = getattr(client, attribute)
		if value not in value_counts:
			value_counts[value] = 1
		else:
			value_counts[value] += 1

	# calculate the entropy of the attribute
	entropy = 0
	total_clients = len(clients)
	# get the value with the highest count, and calculate the probability of that value
	for value, count in value_counts.items():
		probability = count / total_clients
		entropy -= probability * math.log2(probability)


	# so the entropy is defined by how homogenous is the attribute amoung
	# a subset of occurrences
	# the more homogenous, the lower the entropy
	#and subsequently, the more that attribute says about the target attribute
	return entropy

def main():
	clients = [Client() for _ in range(1000)]
	
	canceled_clients = [client for client in clients if client.has_canceled]
	non_canceled_clients = [client for client in clients if not client.has_canceled]
	print(f'Total clients: {len(clients)}')
	print(f'Canceled clients: {len(canceled_clients)}')
	
	attributes = ['age', 'income', 'education']


	for attribute in attributes:
		entropy = calculate_entropy_for_attribute(canceled_clients, attribute)
		print(f'Entropy for cancelled clients {attribute}: {entropy:.4f}')

	for attribute in attributes:
		entropy = calculate_entropy_for_attribute(non_canceled_clients, attribute)
		print(f'Entropy for non-cancelled clients {attribute}: {entropy:.4f}')
	

main()

# at the end of the experiment, i have come to understand that the entropy model